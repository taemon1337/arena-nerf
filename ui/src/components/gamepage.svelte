<script>
  import { page } from '$app/stores'
  import { onMount } from 'svelte'
  import { currentGame, gamelist, pollGame, pollGames, fetchGame } from '$lib/api'
  import { Heading, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Badge, Indicator, Button, ButtonGroup, GradientButton } from 'flowbite-svelte';
  import { ArrowRightOutline } from 'flowbite-svelte-icons';

  export let uuid;

  onMount(() => {
    if (uuid == "current") {
      // nothing changes except current games
      pollGame(uuid)
      pollGames()
    } else {
      fetchGame(uuid)
    }
  })
</script>

<div class="grid grid-rows-3 grid-cols-12 gap-4">
  <div class="row-span-3 col-span-3">
    <Heading tag="h3" class="dark:text-gray-400">Game Teams</Heading>
    {#each Object.entries($currentGame.scoreboard || {}) as [team, count]}
      <Badge class="relative m-4 p-4" color={team}>
        {team}
        <Indicator color="{team}" border size="xl" placement="top-right">
          <span class="text-white text-xs font-bold">{count}</span>
        </Indicator>
      </Badge>
    {/each}

    <Heading tag="h3" class="dark:text-gray-400">Game Nodes</Heading>
    {#each Object.entries($currentGame.nodeboard || {}) as [node, count]}
      <Badge class="p-4 m-4 relative">
        {node}
        <Indicator border size="xl" placement="top-right">
          <span class="text-white text-xs font-bold">{count}</span>
        </Indicator>
      </Badge>
    {/each}

    <p>Started {$currentGame.start_at}</p>
    {#if $currentGame.completed}
    <p>Ended {$currentGame.end_at}</p>
    {:else}
    <p>In progress</p>
    {/if}

    <Table noborder={true} shadow striped={true}>
      <TableHead>
        <TableHeadCell>Game History</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each $gamelist as game}
          <TableBodyRow>
            <TableBodyCell tdClass="px-6 py-0 font-small">
              <a href="/games/{game}">{game}</a>
            </TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>

  <div class="col-span-6 grid-rows-3">
    <Heading tag="h1" class="mb-4 dark:text-gray-400" customSize="text-4xl font-extrabold  md:text-5xl lg:text-6xl">
      {#if uuid == "current"}
        {uuid}
      {/if}
      Game Stats
      <Badge>{uuid}</Badge>
      <Badge color="green">{$currentGame.status}</Badge>
    </Heading>
    <Table striped={true}>
      <TableHead>
        <TableHeadCell></TableHeadCell>
        <TableHeadCell>Rank</TableHeadCell>
        <TableHeadCell>Team name</TableHeadCell>
        <TableHeadCell>Points</TableHeadCell>
      </TableHead>
      <TableBody class="divide-y">
        {#each Object.entries($currentGame.scoreboard || {}) as [team, count]}
        <TableBodyRow>
          <TableBodyCell><Indicator color={team} /></TableBodyCell>
          <TableBodyCell>{team}</TableBodyCell>
          <TableBodyCell>{count}</TableBodyCell>
          <TableBodyCell></TableBodyCell>
        </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
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
    <Table noborder={true} shadow striped={true}>
      <TableHead>
        <TableHeadCell>Game Event Stream</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each [...($currentGame.events || [])].reverse() as event}
          <TableBodyRow>
            <TableBodyCell tdClass="px-6 py-0 font-small">{event}</TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
</div>
