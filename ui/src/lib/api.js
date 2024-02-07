import { writable } from 'svelte/store'

export const API = import.meta.env.VITE_API || "/api/v1"

export const api = (path) => {
  return fetch(API + path, { credentials: 'include' })
}

export const currentStats = writable({
  start_at: "",
  end_at: "",
  length: "5m",
  completed: false,
  status: "game:init",
  teams: [],
  nodes: [],
  events: [],
  nodeboard: {},
  scoreboard: {},
})

export async function fetchStats(id) {
  const res = await api("/" + id + "/stats")
  let data = await res.json()
  if (data.stats) {
    currentStats.update(() => data.stats)
  } else {
    console.error("Unexpected API response, expected a 'stats' key.", data)
  }
}

let poller;

export const pollStats = function(id) {
  if (poller) {
    clearInterval(poller)
  }
  poller = setInterval(function() {
    fetchStats(id)
  }, 5000)
}
