import Vue from 'vue'
import VueRouter from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import LogEntry from './components/LogEntry.vue'
import CreateContest from "./components/CreateContest.vue"
import SettingsView from "./components/SettingsView.vue"

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
    {
        path: '/settings',
        name: 'Settings',
        component: SettingsView
    },
]

export default new VueRouter({
    mode: 'history',
    routes: routes
})