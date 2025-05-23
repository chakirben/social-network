import { handleClientScriptLoad } from "next/script";
import { WS } from "../login/page";

export default async function InitWs() {
    console.log("Initializing WebSocket connection...");
    
    const deferred = new Deferred()
    const ws = new WebSocket("ws://localhost:8080/api/ws");
    ws.onopen = () => {
        console.log("WebSocket connection opened");
        deferred.resolve(ws);
    };
    ws.onmessage = (event) => {
      console.log(event);
    //  handemsj(event.data);
      console.log("Message from server: ", event.data);
    };
    ws.onclose = () => {
        deferred.resolve();
        console.log("WebSocket connection closed");
    };
    ws.onerror = (error) => {
        console.log("WebSocket error: ", error);
        deferred.reject(error);
    };
    await deferred.promise;

}
export class Deferred {
  constructor() {
    this.promise = new Promise((res, rej) => {
    this.resolve = res;
    this.reject = rej;
    });
  }
}
