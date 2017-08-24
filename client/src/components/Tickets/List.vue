<template>
    <div id="list-wrapper" v-if="tickets.list">
    <ul id="ticket-list">
        <li v-for="ticket in tickets">
            {{ ticket.key }}
        </li>
    </ul>
</div>
    <h1 v-else>No tickets!</h1>
</template>
