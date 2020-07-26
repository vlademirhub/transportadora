import * as React from 'react';
import {useState} from "react";
import {useEffect} from "react";
import {Loader} from 'google-maps';
import {useParams} from 'react-router-dom';
import axios from 'axios';

import io from 'socket.io-client';
import {Box} from "@material-ui/core";
import MapControl from "./MapControl";
import OrderInformation from "./OrderInformation";
import {useSnackbar} from "notistack";

const loader = new Loader(process.env.REACT_APP_GOOGLE_API_KEY);
const socket = io(process.env.REACT_APP_MICRO_MAPPING_URL);

const Mapping = () => {
    const {id} = useParams();

    // Hooks
    const [order, setOrder] = useState();
    const [map, setMap] = useState();
    const [startMarker, setStartMarker] = useState();
    const [endMarker, setEndMarker] = useState();
    const [position, setPosition] = useState();

    const snackbar = useSnackbar();// para poder mostrar as mensagem

    useEffect(() => {
        async function load() {
            // Utilizo o axios para consultar o next de pedido e bate no order.controller.ts
            const {data} = await axios
                .get(`${process.env.REACT_APP_MICRO_MAPPING_URL}/orders/${id}`);

            setOrder(data);// jogo para o hook de order

            const [lat, lng] = data.location_geo;
            const position = {lat: parseFloat(lat), lng: parseFloat(lng)}; // guardo as posições como objeto para facilitar o gerenciament

            // jogo dentro de window.google para ficar global para uso, facilita depois.
            window.google = await loader.load();// demora um tempo desconhecido para carregar o google

            const map = new window.google.maps.Map(document.getElementById('map'), {
                center: position,
                zoom: 15,
            });

            // Agora eu crio o marcador de inicio
            const start = new window.google.maps.Marker({
                title: 'Início',
                icon: 'http://maps.google.com/mapfiles/kml/pal4/icon7.png'
            });

            // Agora eu crio o marcador de fim
            const end = new window.google.maps.Marker({
                position: position,
                map: map,// preciso ligar o marcador/ponteiro no mapa
                title: 'Destino'
            });

            // agora eu jogo tudo dentro do meu state através do hook
            setMap(map);
            setStartMarker(start);
            setEndMarker(end);
        }

        load();
    }, [id]);

    // esse aqui permite fazer o watch das variaveis que vai se modificando graças ao hooks
    useEffect(() => {
        // faço a criação da conexão com websocket como client, usando o canal criado pelo back no serviço mapping.service.ts com id do pedido/ordem
        socket.on(`order.${id}.new-position`, (payload) => {
            // conforme eu sou recebendo os dados eu vou atualizando o meu position hook
            console.log(payload);// para ver as comunicações chegando via websocket
            setPosition(payload)
        });
    }, [id]);

    // esse aqui permite fazer o watch das variaveis que vai se modificando graças ao hooks
    useEffect(() => {
        // esse useEffect aqui é o responsavel por atualizar o mapa conforme que o driver/motorista vai andando ele vai atualizando a distancia.

        // como inicialmente eu não vou ter o mapa eu preciso fazer essa verificação abaixo pois são states dependentes.
        if (!map || !position) {
            return;// retorno para não continuar o processamento e gerar erro.
        }

        // Se chegou ao destino, eu mostro a mensagem no rodapé a direita, que chegou ao destino
        if(position.lat === 0 && position.lng === 0){
            snackbar.enqueueSnackbar('Motorista chegou no destino', {
                variant: 'success',
                anchorOrigin: {
                    horizontal: 'right',
                    vertical: 'bottom'
                },
            });
            return;
        }

        // Se não for 0,0 latitude e longitude, eu pego o marcador inicial e posiciono no meu mapa!
        startMarker.setPosition({lat: position.lat, lng: position.lng});
        startMarker.setMap(map);
        const bounds = new window.google.maps.LatLngBounds();// isso aqui é para poder CENTRALIZAR o mapa

        // pego o marcador inicial que eu atualizei em cima com setposition
        // pego o marcador final que eu atualizei em cima com setposition
        bounds.extend(startMarker.getPosition());
        bounds.extend(endMarker.getPosition());

        map.fitBounds(bounds);// ajusto aqui o status do trajeto centralizando o mapa.
    }, [map, position, snackbar, startMarker, endMarker]);

    return (// template do mapa
        <div id={'map'} style={{width: '100%', height: '100%'}}>
            {
                map && // só mostra se tiver o mapa!
                <MapControl map={map} position={window.google.maps.ControlPosition.TOP_RIGHT}>
                    <Box m={'10px'}>
                        <OrderInformation order={order}/>
                    </Box>
                </MapControl>
            }
        </div>
    );
};

export default Mapping;
