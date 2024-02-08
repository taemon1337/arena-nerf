<script>
  import { page } from '$app/stores'
  import GameEventsTable from '$src/components/game-events-table.svelte'
  import GameHistoryTable from '$src/components/game-history-table.svelte'
  import GameTeams from '$src/components/game-teams.svelte'
  import GameNodes from '$src/components/game-nodes.svelte'
  import Scoreboard from '$src/components/scoreboard.svelte'
  import { onMount } from 'svelte'
  import { currentGame, gamelist, pollGame, pollGames, fetchGame } from '$lib/api'
  import { List, Li, Heading, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Badge, Indicator, Button, ButtonGroup, GradientButton } from 'flowbite-svelte';
  import { ArrowRightOutline, CheckCircleSolid } from 'flowbite-svelte-icons';

  export let uuid

  onMount(() => {
    if (uuid == "current") {
      // nothing changes except current games
      pollGame(uuid)
      pollGames()
    } else {
      fetchGame(uuid)
      fetchGame("all")
    }
  })
</script>

<div class="grid grid-rows-3 grid-cols-12 gap-4">
  <div class="row-span-3 col-span-3">
    <GameTeams />
    <GameNodes />

    <p>Started {$currentGame.start_at}</p>
    {#if $currentGame.completed}
    <p>Ended {$currentGame.end_at}</p>
    {:else}
    <p>In progress</p>
    {/if}

    <GameHistoryTable />
  </div>

  <div class="col-span-6 grid-rows-3">
    <Scoreboard uuid={uuid} />
    <ButtonGroup class="space-x-px">
      {#if $currentGame.completed}
      <GradientButton color="purpleToBlue">New Game</GradientButton>
      {:else}
      <GradientButton color="purpleToBlue">Stop Game</GradientButton>
      {/if}
      <GradientButton color="cyanToBlue">Settings</GradientButton>
    </ButtonGroup>
  </div>
  <div class="row-span-2 col-span-3">
    <GameEventsTable />
  </div>
</div>
