export const API = import.meta.env.VITE_API || "/api/v1"

export const api = (path) => {
  return fetch(API + path, { credentials: 'include' })
}

export const currentApi = {
  get: async () => {
    const res = await api("current/stats");
		let stats = await res.json();
    return stats.stats;
  }
}
