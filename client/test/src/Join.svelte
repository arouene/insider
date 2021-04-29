<script>
 import { onDestroy } from 'svelte';
 import { game } from './store.js';
 import { request } from './utils.js';

 let salon;
 let name;

 function join() {
     let path = {game: salon, path: "join"};
     let args = new Map([["playerName", name]]);
     request(path, args, (j) => { game.set(salon); });
 }
</script>

<label for="salon">Salon</label>
<input bind:value={salon} id="salon" type="text"/>
<label for="name">Nom</label>
<input bind:value={name} id="name" type="text"/>
<button on:click={join}>Join</button>

<style>
</style>
