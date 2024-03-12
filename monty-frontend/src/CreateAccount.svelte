<script>
    const handleSubmit = async e => {
        let ACTION_URL = e.target.action;
        let formData = new FormData(e.target);
        let jsonData = {};
        formData.forEach((value, key) => {jsonData[key] = value});

        let res = await fetch(ACTION_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(jsonData)
        })

        if (res.status === 200) {
            let data = await res.json();
            console.log(data);
        } else {
            console.log("Error");
        }
    }
</script>

<h1>Create Account</h1>
<form action="api/create-account" on:submit|preventDefault={handleSubmit}>
    <label for="username">Username:</label>
    <input type="text" id="username" name="username">
    <label for="password">Password:</label>
    <input type="password" id="password" name="password">
    <label for="email">Email:</label>
    <input type="email" id="email" name="email">
    <input type="submit" value="Create Account">
</form>

<style>
    form {
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    input {
        margin: 10px;
    }
    input[type="submit"] {
        margin: 10px;
    }
    label {
        margin: 10px;
    }
    h1 {
        margin: 10px;
    }
</style>