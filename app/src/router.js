import Vue from 'vue'
import VueRouter from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import LogEntry from './components/LogEntry.vue'
import LogTest from './components/LogTest.vue'

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
        path: '/logs2',
        name: 'Logs',
        component: LogTest
    }
]

export default new VueRouter({
    mode: 'history',
    routes: routes
})