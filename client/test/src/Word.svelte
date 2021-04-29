<script>
 import { state, game, host } from './store.js';
 import { request } from './utils.js';

 export let canEdit;

 function newWord() {
     request("word");
 }
</script>

{#if $state.word}
<h2>Votre mot est «&nbsp;<span id="word">{$state.word}</span>&nbsp;»</h2>
{#if canEdit}
<button on:click={newWord}>Nouveau mot</button>
{/if}
{/if}

<style>
 #word {
     font-weight: bolder;
 }
</style>
