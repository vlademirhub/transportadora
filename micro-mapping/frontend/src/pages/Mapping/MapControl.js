import * as React from 'react';
import {useMemo} from "react";
import {createPortal} from "react-dom";
import {useEffect} from "react";


// A gente precisa criar utilizando o conceito do react chamado de PORTAL porque a gente precisa criar um elemento que não está ligado a DIV principal do template.
const MapControl = (props) => {
    const {map, position, children} = props;
    const controlDiv = useMemo(() => document.createElement('div'), []);

    useEffect(() => {
        if (map && position) {
            map.controls[position].push(controlDiv);
        }
    }, [map, position, controlDiv]);

    return createPortal(children, controlDiv);
};


export default MapControl;
