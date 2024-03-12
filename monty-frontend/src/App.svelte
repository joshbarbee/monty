<script>
  import { navOptions } from "./Navbar.svelte";
  let selected = navOptions[0]
  let intSelected = 0

  function changeComponent(event) {
    selected = navOptions[event.srcElement.id]
    intSelected = event.srcElement.id
  }
</script>

<main>
  <nav>
    <ul> 
      {#each navOptions as option, i}
        <li>
          <button on:click={changeComponent} id={i.toString()} role="tab">{option.page}</button>
        </li>
      {/each}
    </ul>
  </nav>
  <div class='main-content'>
    <svelte:component this={selected.component}/>
  </div>
</main>

<style>
  main {
    display: block;
    flex-direction: column;
    align-items: center;
    width: 100%;
    height: 100%;
  }

  nav {
    position: sticky;
    top: 0;
  }

  ul {
    display: flex;
    list-style: none;
    padding: 0;
  }

  li {
    margin: 0 10px;
  }

  button {
    padding: 10px;
  }

  [role="tab"] {
    background-color: #f0f0f0;
    color: #000;
    border: 1px solid #ccc;
    border-radius: 5px;
  }

  [role="tab"]:hover {
    background-color: #e0e0e0;
  }

  .main-content {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
  }
</style>
