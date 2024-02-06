<script>
  import { onMount } from 'svelte'
  import { currentApi } from '$lib/api'
  import { Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Checkbox } from 'flowbite-svelte';
  import { LinkOutline } from 'flowbite-svelte-icons';
  let routes = []
  let stats = {}

  onMount(async () => {
		stats = await currentApi.get() || [];
	});

</script>

<div class="p-2 m-2">
  <h3>Current Game Stats</h3>
  <Table hoverable={true}>
    <TableHead>
      <TableHeadCell class="!p-4">
        <Checkbox />
      </TableHeadCell>
      <TableHeadCell>Name</TableHeadCell>
      <TableHeadCell>User</TableHeadCell>
      <TableHeadCell>Link</TableHeadCell>
      <TableHeadCell>
        <span class="sr-only">Approve</span>
      </TableHeadCell>
    </TableHead>
    <TableBody class="divide-y">
      {#each routes as route}
      <TableBodyRow>
        <TableBodyCell class="!p-4">
          <Checkbox />
        </TableBodyCell>
        <TableBodyCell>{route.metadata.name}</TableBodyCell>
        <TableBodyCell>{route.spec.user}</TableBodyCell>
        <TableBodyCell>
          <a href="{KACM_PROXY}/users/{route.spec.user}/{route.spec.access_path}/" target="_blank">
            <LinkOutline />
          </a>
        </TableBodyCell>
        <TableBodyCell>
          <a href="/tables" class="m-1 font-medium text-primary-600 hover:underline dark:text-primary-500">Approve</a>
          <a href="/tables" class="m-1 font-medium text-primary-600 hover:underline dark:text-primary-500">Disable</a>
        </TableBodyCell>
      </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
</div>
