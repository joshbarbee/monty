<script context='module'>
    import Login from "./Login.svelte";
    import Home from "./Home.svelte";
    import Logout from "./Logout.svelte";
    import CreateAccount from "./CreateAccount.svelte";

    export const navOptions = [
        { page: 'Home', component : Home},
        { page: 'Login', component : Login},
        { page: 'Logout', component : Logout},
        { page: 'Create Account', component : CreateAccount}
    ]
</script>