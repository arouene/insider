<script>
 export let players = [];
</script>

<div id="lobby">
    <h2>Joueurs:</h2>
    <ul>
{#each players as player (player.id)}
        <li>{player.name}</li>
{/each}
    </ul>
</div>
