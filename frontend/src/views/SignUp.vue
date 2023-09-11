<script setup>

import { ref } from 'vue';
import router from '../router';

const name = ref("")
const email = ref("")
const password = ref("")

async function submitForm(){
    try {
        const req = await fetch("http://localhost:8000/signup", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                name: name.value,
                email : email.value,
                password: password.value,
            })
        })

        if (req.status === 200){
            const res = await req.json()
            console.log(res)
            router.replace("/login")
        } else {
            console.log("Error occured")
            resetCredentials()
        }

    } catch (error) {
        console.error(error)
    }

}

function resetCredentials() {
    name.value = ""
    email.value = ""
    password.value = ""
}

</script>

<template>
    <h3>SignUp page</h3>
    <form id="signupform">
        <label for="name">Name</label>
        <input type="text" id="name" placeholder="Name" v-model="name">
        <label for="email">Email</label>
        <input type="email" id="email" placeholder="Email" v-model="email">
        <label for="password">Password</label>
        <input type="password" id="password" placeholder="password" v-model="password">
    </form>
    <button v-on:click="submitForm">Submit</button>

</template>