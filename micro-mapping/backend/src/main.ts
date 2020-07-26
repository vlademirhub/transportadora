import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import {NestExpressApplication} from "@nestjs/platform-express";

async function bootstrap() {
  const app = await NestFactory.create<NestExpressApplication>(AppModule, {cors: true});
  // Preciso ativar esse cors por causa da politica de troca cruzada
  // se na hora de dar resposta, o nest não informar que o cors está ativo o BROWSER vai bloquear o React de acessar as informações.


  await app.listen(3001);
}
bootstrap();

// http://localhost:3001

// http://localhost:3002 - mapa frontend
