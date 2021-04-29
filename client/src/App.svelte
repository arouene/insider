<script>
 import { onMount } from 'svelte';
 import Lobby from './Lobby.svelte';
 import Word from './Word.svelte';

const apiUrl = "http://localhost:3000/game/";

 let joinning = true;
 let fields = {
     "teamName": "",
     "playerName": "",
     "uuid": "",
     "role": "",
     "players": [],
     "word": "",
     "state": "",
 }
 let master = true;
 let insider = false;

 var timeout;
 onMount(
     async () => {
         setInterval(() => {
             console.log("Refresh state");
             let res = fetch(apiUrl + fields.teamName + "/state?uuid=" + fields.uuid).then(r => {
                 if (r.ok) {
                     return r.json();
                 }
             })
            .then(json => {
                console.log(json);
                fields.players = [];
                for (const player in json.players) {
                    fields.players.push(json.players[player]);
                }
                console.log(fields.players);
            })
            .catch(err => {});
         }, 3000);
     }
 );

 function getState() {
     return function() {
     };
 }

async function joinTeam() {
    var team = document.getElementById("team");
    var btn = team.getElementsByTagName("button");
    btn[0].disabled = true;

    const res = await fetch(apiUrl + fields.teamName + "?playerName=" + fields.playerName, {});
    const player = await res.json();

    if (res.ok) {
        fields.uuid = player.id
        fields.role = player.role
        btn[0].disabled = false
        team.style = "display: none;"
    } else {
        btn[0].disabled = false
    }
 }

function debounce(func, wait) {
    var timeout;
    return function() {
        var context = this, args = arguments;
        var later = function() {
            timeout = null;
            func.apply(context, args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
};

const cryptPlaintext = debounce(event => {
    genService("crypt", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({secret: event.target.value})
    }).apply();
}, 300)

async function start() {
    var btn = document.getElementById("start");
    btn.disabled = true;

    const res = await fetch(apiUrl + fields.teamName + "?uuid=" + fields.uuid, {});
    const player = await res.json();

    if (res.ok) {
        fields.uuid = player.id
        fields.role = player.role
        btn[0].disabled = false
        team.style = "display: none;"
    } else {
        btn.disabled = false;
    }
}

async function stop() {
    var btn = team.getElementById("stop");
    btn.disabled = true;

    const res = await fetch(apiUrl + fields.teamName + "?uuid=" + fields.uuid, {});
    const player = await res.json();

    if (res.ok) {
        fields.uuid = player.id
        fields.role = player.role
        btn[0].disabled = false
        team.style = "display: none;"
        getPlayers().apply()
    } else {
        btn.disabled = false;
    }
}
</script>

<style>
 :global(body) {
     background: circular-gradient(#f22 29.25%,#f60 100%) no-repeat fixed;
     background-color: #f90;
 }

 h1 {
     color: white;
     text-transform: uppercase;
     text-align: center;
     font-size: 4em;
     font-weight: 100;
 }

 input {
     width: 100%;
     font-size: 0.95em;
 }
</style>

<h1>Insider</h1>

{#if joinning}
<div id="team">
    <label>
        <h2>Equipe</h2>
        <input bind:value={fields.teamName}>
        <h2>Joueur</h2>
        <input bind:value={fields.playerName}>
    </label>
    <button on:click={joinTeam}>Joindre</button>
</div>
{/if}

<Lobby players={fields.players}/>
{#if master || insider}
<Word/>
{/if}

<div id="stop">
    <button on:click={stop}>Stop</button>
</div>
