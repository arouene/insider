import { onDestroy } from 'svelte';
import { state, uuid, game } from './store.js';

const schema = "http://"
const host = "localhost:3000";

let id, party;
uuid.subscribe(value => { id = value; });
game.subscribe(value => { party = value; });

export function request(path, args=[], action=(j)=>{}) {
    let url_path = "/game/";
    // Construct url from the function parameters
    // if the path is a string, it optionaly just contains the rest
    // of the path after the game party name, but if the path
    // is an object, it can contains the gane party name.
    // obect is : {game: "value", path: "value"}
    if (typeof path === "string") {
        url_path += party + (path ? "/" + path : "");
    } else if (typeof path === "object") {
        if (path.game) {
            url_path += path.game;
        } else {
            url_path += party;
        }
        if (path.path) {
            url_path += "/" + path.path;
        }
    }

    // let construct options string with the function parameters
    // if args is valued, each key, value will be translated as an
    // argument string "&key=value", uuid is always added, as it's
    // mandatory for the backend.
    let url_args = '?uuid=' + id;
    args.forEach((v, k) => {
        url_args += "&" + k + "=" + v;
    });

    const url = schema + host + url_path + url_args;

    fetch(url)
        .then(res => {
            if (res.ok) {
                return res.json();
            } else {
                console.log("Mauvaise reponse du serveur");
            }
        })
        .then(json => {
            console.log(json);
            return action(json);
        })
        .catch(error => {
            console.log("Il y a eu un probleme avec la requete: " + error.message);
        });
}
