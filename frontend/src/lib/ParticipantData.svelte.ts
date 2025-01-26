import { error } from "@sveltejs/kit"
import type { Participants } from "./types"

export class ParticipantsData {
    participant: Participants = $state([])
    eventID: string

    private constructor(eventID: string, data: Participants) {
        this.participant = data
        this.eventID = eventID
    }

    get(): Participants {
        return this.participant
    }

    static async init(eventID: string): Promise<ParticipantsData> {
        const response = await fetch(`/api/event/${eventID}/participant`)
        if (response.ok) {
            const data: Participants = await response.json()
            return new ParticipantsData(eventID, data)
        } else {
            error(500, "Failed to get a list of participants")
        }
    }

    async delete(index: number) {
        const toDelete = this.participant[index]

        // delete locally to update the UI
        this.participant = this.participant.filter((_,i) => i != index)

        // delete on the server
        if (toDelete.id == '') {
            // it's a row that was never actually saved on the server, so we can just ignore it
            return
        }

        const response = await fetch(`/api/event/${this.eventID}/participant/${toDelete.id}`, {
            method: 'DELETE'
        })
        if (!response.ok) {
            console.error(response)
            error(500, "Nepodařilo se uložit odstranění účastníka.")
        }
    }

    async addNew() {
        // add locally to trigger UI update
        const newIndex = this.participant.length
        this.participant = [...this.participant, {
            id: '',
            email: '',
            name: '',
        }]

        // add on the server
        const response = await fetch(`/api/event/${this.eventID}/participant`, {
            method: 'POST',
            body: JSON.stringify(this.participant[newIndex])
        })
        if (!response.ok) {
            //this.delete(newIndex)
            console.log(response)
            error(500, "Nepodařilo se přidat účastníka.")
        } else {
            const data = await response.json()
            this.participant[newIndex].id = data.id
        }
    }

    async notifyUpdate(index: number) {
        if (this.participant[index].id == '') {
            // it's a row that was never actually saved on the server, so we can just ignore it
            alert("This row was never saved on the server, this should not happen.")
            return
        }

        const response = await fetch(`/api/event/${this.eventID}/participant/${this.participant[index].id}`, {
            method: 'POST',
            body: JSON.stringify(this.participant[index])
        })

        if (!response.ok) {
            console.error(response)
            error(500, "failed saving the data")
        }
    }
}