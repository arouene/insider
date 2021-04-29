import { readable, writable } from 'svelte/store';
import { v4 as uuidv4 } from 'uuid';

export const host = "localhost:3000";

const refreshInterval = 2000;
const idCookieName = "uuid";

// Store for game name (play team)
export const game = writable(undefined);

let gameName;
game.subscribe(value => { gameName = value; });

// Store for player uuid, must remain the same even with page refresh
export const uuid = readable(undefined, function start(set) {
    if (sessionStorage.getItem(idCookieName)) {
        // Get id from cookie
        set(sessionStorage.getItem(idCookieName));
    } else {
        // or set a new id
        let id = uuidv4();
        sessionStorage.setItem(idCookieName, id);
        set(id);
    }
});

let id;
uuid.subscribe(value => { id = value; });

// Store for game state
export const state = readable({players: []}, function start(set) {
    const interval = setInterval(() => {
        if (gameName == undefined) {
            return;
        }
        fetch('http://' + host + '/game/' + gameName + '/state?uuid=' + id)
            .then(res => {
                if (res.ok) {
                    return res.json();
                } else {
                    console.log("Mauvaise reponse du serveur");
                }
            })
            .then(json => {
                json.state.players.sort(sortPlayers);
                set(json.state);
            })
            .catch(error => {
                console.log("Il y a eu un probleme avec la requete: " + error.message);
            });
    }, refreshInterval);

    return function stop() {
        clearInterval(interval);
    };
});

function sortPlayers(a, b) {
    // Current player is first
    if (a.id.length > 0) {
        return -1;
    } else if (b.id.length > 0) {
        return 1;
    }

    // Then sorted by name
    let nameA = a.name.toLowerCase();
    let nameB = b.name.toLowerCase();

    if (nameA < nameB) {
        return -1;
    } else if (nameA > nameB) {
        return 1;
    } else {
        return 0;
    }
}

export function getPlayer(state) {
    return state.players.find(v => v.id === id);
}
