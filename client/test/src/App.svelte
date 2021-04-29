<script>
 import { state, game, getPlayer } from './store.js';
 import { request } from './utils.js';
 import Word from './Word.svelte';
 import Join from './Join.svelte';
 import Lobby from './Lobby.svelte';

 function ready() {
     request("ready");
 }

 function start() {
     request("start");
 }

 function reset() {
     request("stop");
     request("reset");
 }

 $: readyDisabled = $state.players.length < 4;
 $: player = getPlayer($state);
 $: phase = $state.phase;
 $: canEdit = player ? player.role === "MASTER" : false;
</script>

<main>
	<img alt="logo" src="../logo.svg" width="100px">
	<h1>Insider</h1>
	<img alt="logo" src="../logo.svg" width="100px">
</main>

<Lobby/>

<Word canEdit={canEdit}/>

{#if phase == undefined}
	<Join/>
{:else if phase === "CREATED"}
	<button on:click={ready} disabled={readyDisabled}>Ready</button>
{:else if phase === "SETUP" || phase === "STARTED"}
    Vous êtes un {player.role}
{/if}
{#if phase === "STARTED"}
    <button on:click={reset}>Reset</button>
{/if}

<style>
 :global(body) {
     background-color: #f70003;
 }

 :global(button) {
     background-color: #f70003;
     color: #000;
     border: 1px solid black;
     transition-duration: 0.4s;
 }

 :global(button:hover) {
     background-color: #000;
     color: #f70003;
 }

 main {
     text-align: center;
     padding: 1em;
     max-width: 240px;
     margin: 0 auto;
 }

 h1 {
     color: black;
     text-transform: uppercase;
     font-size: 4em;
     font-weight: 400;
     margin-left: 1em;
     margin-right: 1em;
 }

 @media (min-width: 640px) {
     main {
         max-width: none;
     }
}
</style>
