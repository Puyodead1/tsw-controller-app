import { create } from 'zustand'
import { EventsOn } from '../wailsjs/runtime/runtime'

type StoreType = {
  logs: string[]
}

export const useStore = create<StoreType>((set) => ({
  logs: [],
}))

EventsOn('log', (log: string) => {
  useStore.setState(({ logs }) => ({ logs: [...logs, log] }))
})