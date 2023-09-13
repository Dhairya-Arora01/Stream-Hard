import { createStore } from 'vuex'

const store = createStore({
    state(){
        return {
            err: ""
        }
    },
    mutations: {
        addErr(state, msg){
            state.err = msg
        },
        clearErr(state){
            state.err = ""
        },
        funcErr(state, msg){

            state.err = msg
            
            setTimeout(()=>{
                state.err=""
            }, 5000)
        }
    },
    getters: {
        getErr(state){
            return state.err
        }
    }
})

export default store