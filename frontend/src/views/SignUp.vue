<script setup>

import { ref } from 'vue';
import router from '../router';

const name = ref("")
const email = ref("")
const password = ref("")

async function submitForm(){
    try {
        const req = await fetch("http://192.168.49.2:31768/signup", {
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
            store.commit("funcErr", "User already exists")
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
    <div id="signupform">
        <h3>SignUp</h3>
        <div id="form">
            <div id="formelement">
                <label for="name">Name</label>
                <input type="text" id="name" placeholder="Name" v-model="name">
            </div>
            <div id="formelement">
                <label for="email">Email</label>
                <input type="email" id="email" placeholder="Email" v-model="email">
            </div>
            <div id="formelement">
                <label for="password">Password</label>
                <input type="password" id="password" placeholder="password" v-model="password">
            </div>
        </div>
        <button v-on:click="submitForm">SignUp</button>
        <p>Already a member? <router-link to="login">Login</router-link></p>
    </div>

</template>

<style>

@import url('https://fonts.googleapis.com/css2?family=Lato&family=Poppins&display=swap');

div#signupform {
    display: flex;
    flex-direction: column;
    justify-content:space-around;
    align-items: center;
    background-color: rgb(127, 137, 150);
    height: 70%;
    width: 30%;
    font-family: 'Poppins', sans-serif;
    font-weight: bold;
    border-radius: 0.5em;
}

div#signupform > h3 {
    font-size: 1.7em;
}

div#form {
    width: 80%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
}

div#formelement {
    display: flex;
    flex-direction: column;
    margin: 2%;
    width: 70%;
}

div#formelement > input {
    height: 2.8em;
    font-size: 0.9em;
    color: white;
    border-radius: 0.2em;
    border: none;
    padding-left: 0.4em;
    background-color: rgb(99, 106, 116);
}

div#signupform > button {
    width: 20%;
    height: 7%;
    background-color: rgb(99, 106, 116);
    font-size: 1.05em;
    color: white;
    border: none;
}

div#signupform > button:hover {
    background-color: rgb(48, 51, 55);
}

</style>