import { create } from 'zustand'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { events } from './events'

type StoreType = {
  logs: string[]
}

export const useStore = create<StoreType>((set) => ({
  logs: [],
}))

EventsOn(events.log, (log: string) => {
  useStore.setState(({ logs }) => ({ logs: [...logs, log] }))
})