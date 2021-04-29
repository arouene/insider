<script>
 let master = true;
 let insider = false;
 let word = "Hello word"
 let setup = true;

 function newWord() {

 }

 function start() {

 }
</script>

<h1>Mot: {word}</h1>

{#if master && setup}
<button on:click={newWord}>Nouveau mot</button>
<button on:click={start}>Commencer</button>
{/if}
