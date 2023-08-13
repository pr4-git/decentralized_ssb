import {  createMemo,createSignal, createResource,  onMount } from "solid-js"
import { ViewOneProfile, GetMyPublickey, ViewAllProfiles } from "../wailsjs/go/main/App"
import { minidenticon } from 'minidenticons';

type Profile = { pubkey: string }

function Image({ pubkey }: Profile) {

    const svgURI = createMemo(() => "data:image/svg+xml;utf8," + encodeURIComponent(minidenticon(pubkey)))

    return (
        <div class="w-24 h-24 border rounded-full flex items-center justify-center">
            <img height={80} width={80} src={svgURI()} alt="profile picture" />
            <div ></div>
        </div>
    )
}

const [updatepfcount, setUpdatePfCount] = createSignal(0)
export const UpdateProfile = () => setUpdatePfCount(updatepfcount()+1)
const [profiles, { mutate, refetch }] = createResource(updatepfcount,ViewAllProfiles)
const [selfKey] = createResource(GetMyPublickey)

export function ProfileViewer({ pubkey }: Profile) {

    const getProfile = () => profiles()?.filter(profile => profile.pubkey.toString() == pubkey)[0]

    const profileName = createMemo(() => {
        if (pubkey == selfKey())
            return "Me"
        else if (getProfile() == undefined)
            return "Anonymous"
        else return getProfile()?.alias
    })


    onMount(refetch)

    return (
        <div>
            <Image pubkey={pubkey} />
            <div class="flex flex-col items-center">
                <div>{profileName()}</div>
                <div class="text-xs hover:underline"
                    title="click to copy"
                >@{pubkey.substring(0, 8)}</div>
            </div>
        </div>
    )
}