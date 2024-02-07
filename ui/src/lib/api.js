import { writable } from 'svelte/store'

export const API = import.meta.env.VITE_API || "/api/v1"

export const api = (path) => {
  return fetch(API + path, { credentials: 'include' })
}

export const currentGame = writable({
  start_at: "",
  end_at: "",
  length: "",
  completed: false,
  status: "no active game",
  teams: [],
  nodes: [],
  events: [],
  nodeboard: {},
  scoreboard: {},
})

export const gamelist = writable([])

export async function fetchGame(id) {
  const res = await api("/games/" + id)
  let data = await res.json()
  if (data.stats) {
    currentGame.update(() => data.stats)  // games/:uuid returns {stats:{}}
  } else if (data.games) {
    gamelist.update(() => data.games) // games/all returns {games:[]}
  } else {
    console.error("Unexpected API response, expected a 'stats' key.", data)
  }
}

export async function fetchGames() {
  return fetchGame("all")
}

let poller;
let pollgames;

export const pollGame = function(id) {
  if (poller) {
    clearInterval(poller)
  }
  poller = setInterval(function() {
    fetchGame(id)
  }, 5000)
}

export const pollGames = function() {
  if (pollgames) {
    clearInterval(pollgames)
  }
  pollgames = setInterval(function() {
    fetchGame("all")
  }, 5000)
}
