<script>
	import '../app.postcss';
  import { Avatar, Badge, Navbar, NavBrand, NavLi, NavUl, NavHamburger, DarkMode, Dropdown, DropdownItem, DropdownHeader, DropdownDivider} from 'flowbite-svelte'
  import { UserSolid } from 'flowbite-svelte-icons'

  let btnClass="text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-gray-200 dark:focus:ring-gray-700 rounded-lg text-sm p-2"
</script>

<Navbar let:hidden let:toggle>
  <NavBrand href="/">
    <img
      src="https://flowbite.com/docs/images/logo.svg"
      class="mr-3 h-6 sm:h-9"
      alt="Flowbite Logo"
    />
    <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
      Kacm
    </span>
  </NavBrand>
  <NavHamburger />
  <NavUl>
    <NavLi href="/">Home</NavLi>
    <NavLi>
      <DarkMode title="Toggle dark and light mode" class={btnClass} />
    </NavLi>
    <NavLi>
    </NavLi>
    <div class="flex items-center md:order-2">
      <UserSolid id="avatar-menu" />
      <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
    </div>
    <Dropdown placement="bottom" triggeredBy="#avatar-menu">
      <DropdownHeader>
        <span class="block text-sm">User</span>
      </DropdownHeader>
      <DropdownItem>Dashboard</DropdownItem>
      <DropdownItem>Settings</DropdownItem>
      <DropdownItem>Earnings</DropdownItem>
      <DropdownDivider />
      <DropdownItem>Sign out</DropdownItem>
    </Dropdown>
  </NavUl>
</Navbar>

<div class="container mx-auto bg-white rounded-xl overflow-y-auto h-full">
  <slot />
</div>
