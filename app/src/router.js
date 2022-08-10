import Vue from 'vue'
import VueRouter from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import LogEntry from './components/LogEntry.vue'
import CreateContest from "./components/CreateContest.vue"

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'HelloWorld',
        component: HelloWorld
    },
    {
        path: '/logs',
        name: 'Logs',
        component: LogEntry
    },
    {
        path: '/create_contest',
        name: 'Create Contest',
        component: CreateContest
    },
]

export default new VueRouter({
    mode: 'history',
    routes: routes
})