<script setup>
import { computed, onMounted, ref } from 'vue';
import store from '../store';

const err = computed(()=>store.getters.getErr)

</script>

<template>

    <div id="err" v-if="err !== ''">
        <p>{{ err }}</p>
    </div>

</template>

<style>

div#err{
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    position: fixed;
    top: 15%;
    left: 50%;
    transform: translateX(-50%);
    width: 20%;
    height: 5%;
    text-align: center;
    font-family: 'Poppins', sans-serif;
    color: white;
    background-color: rgb(179, 11, 11);
    border-radius: 2px;
}

</style>