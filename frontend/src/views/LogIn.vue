<script setup>
import { ref } from 'vue';
import router from '../router';

const email = ref("")
const password = ref("")

async function submitForm() {
    try {
        const req = await fetch("http://localhost:8000/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email : email.value,
                password: password.value,
            })
        })

        if (req.status === 200){
            const res = await req.json()
            localStorage.setItem("token", res["token"])
            router.replace("/home")
        } else {
            console.log("Error occured")
            resetCredentials()
        }

    } catch (error) {
        console.error(error)
    }
}

function resetCredentials() {
    email.value = ""
    password.value = ""
}

</script>

<template>
    <h3>Login</h3>
    <form>
        <label for="email">Email</label>
        <input type="email" id="email" placeholder="Email" v-model="email">
        <label for="password">Password</label>
        <input type="password" id="password" placeholder="Password" v-model="password">
    </form>
    <button v-on:click="submitForm">Login</button>
</template>