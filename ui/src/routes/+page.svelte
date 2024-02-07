<script>
  import { onMount } from 'svelte'
  import { currentApi } from '$lib/api'
  import { Heading, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Badge, Indicator } from 'flowbite-svelte';
  import { ArrowRightOutline } from 'flowbite-svelte-icons';

  let default_stats = {
    start_at: "",
    end_at: "",
    length: "5m",
    completed: false,
    status: "game:init",
    teams: [],
    nodes: [],
    events: [],
    nodeboard: {},
    scoreboard: {},
  }
  let stats = default_stats

  onMount(async () => {
		stats = await currentApi.get() || default_stats;
	});

</script>

<div class="grid grid-rows-3 grid-cols-12 gap-4">
  <div class="row-span-3 col-span-3">
    <Heading tag="h3" class="dark:text-gray-400">Current Teams</Heading>
    {#each Object.entries(stats.scoreboard || {}) as [team, count]}
      <Badge class="relative m-4 p-4" color={team}>
        {team}
        <Indicator color="{team}" border size="xl" placement="top-right">
          <span class="text-white text-xs font-bold">{count}</span>
        </Indicator>
      </Badge>
    {/each}

    <Heading tag="h3" class="dark:text-gray-400">Active Game Nodes</Heading>
    {#each Object.entries(stats.nodeboard || {}) as [node, count]}
      <Badge class="p-4 m-4 relative">
        {node}
        <Indicator border size="xl" placement="top-right">
          <span class="text-white text-xs font-bold">{count}</span>
        </Indicator>
      </Badge>
    {/each}

    <p>Started {stats.start_at}</p>
    {#if stats.completed}
    <p>Ended {stats.end_at}</p>
    {:else}
    <p>In progress</p>
    {/if}
  </div>

  <div class="col-span-6 grid-rows-3">
    <Heading tag="h1" class="mb-4 dark:text-gray-400" customSize="text-4xl font-extrabold  md:text-5xl lg:text-6xl">
      Current Game Stats
      <Badge color="green">{stats.status}</Badge>
    </Heading>
    <Table striped={true}>
      <TableHead>
        <TableHeadCell></TableHeadCell>
        <TableHeadCell>Rank</TableHeadCell>
        <TableHeadCell>Team name</TableHeadCell>
        <TableHeadCell>Points</TableHeadCell>
      </TableHead>
      <TableBody class="divide-y">
        {#each Object.entries(stats.scoreboard || {}) as [team, count]}
        <TableBodyRow>
          <TableBodyCell><Indicator color={team} /></TableBodyCell>
          <TableBodyCell>{team}</TableBodyCell>
          <TableBodyCell>{count}</TableBodyCell>
          <TableBodyCell></TableBodyCell>
        </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
  <div class="row-span-2 col-span-3">
    <Table noborder={true} shadow striped={true}>
      <TableHead>
        <TableHeadCell>Game Event Stream</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each [...(stats.events || [])].reverse() as event}
          <TableBodyRow>
            <TableBodyCell tdClass="px-6 py-0 font-small">{event}</TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
</div>
