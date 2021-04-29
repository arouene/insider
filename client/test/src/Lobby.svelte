<script>
 import { state } from './store.js';
</script>

<div class="box">
    <h2>Joueurs&nbsp;:</h2>
    <ul>
    {#each $state.players as player}
        <li>{player.name}
            {player.role !== "NONE" ? "(" + player.role + ")" : ""}
        </li>
    {/each}
    </ul>
</div>

<style>
 .box {
     background-color: #fff;
     box-shadow: 0 0 8px -1px;
     padding: 0.5em 0.75em;
     margin: 1em;
 }
 h2 {
     margin: 0;
 }
</style>
