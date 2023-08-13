import { ViewAllProfiles, AddNewPeerProfile } from "../wailsjs/go/main/App"
import { Show, For, onMount, createResource, createSignal,createMemo } from "solid-js";
import { Icon } from "./Icon";
import { minidenticon } from "minidenticons";
import userimg from "./assets/userImg.svg";

export function Following() {
    return (<div class="w-[100%] grid grid-rows-[30%_5%_65%]">
        <div class="flex justify-center "><AddPeer /></div>
        <div class="w-[max%] drop-shadow shadow-lg shadow-[#242424] flex place-content-center place-items-center">
        </div>
        <div class="flex justify-center"><KnownPeerList /></div>
    </div>)
}

const [fetchCount, setFetchCount] = createSignal(0)
const refetch = () => setFetchCount(fetchCount() + 1)



function AddPeer() {
    const [pubkey, setPubkey] = createSignal("")
    const [alias, setAlias] = createSignal("")
    const [flash, setFlash] = createSignal("")

    const flashNotif = (msg: string) => {
        let before = flash()
        setFlash(msg)
        setTimeout(() => setFlash(before), 1000)
    }

    const CreatePeer = () => {
        AddNewPeerProfile(pubkey(), alias())
            .then(() => setFlash("Success!"))
            .then(refetch)
            .catch(err => setFlash(err))
    }

    return (<div class="flex flex-col place-center gap-4">
        <Icon iconSrc={userimg} msg="Add a new peer"></Icon>
        <div class="flex flex-row gap-2">
            <input class="w-96 h-8 bg-gray-800 text-center border border-gray-700 placeholder:text-white-600"
                placeholder="Public Key"
                onInput={(e) => { setPubkey(e.target.value) }}
            >
            </input>
        </div>
        <div class="flex flex-row gap-2 justify-evenly">
            <input class="w-72 h-8 bg-gray-800 text-center border border-gray-700 placeholder:text-white-600"
                placeholder="Alias"
                onInput={(e) => { setAlias(e.target.value) }}
            >
            </input>
            <button class="bg-blue-500 border border-blue-500 rounded-full px-3 py-1"
                onclick={CreatePeer}
            >
                Add Peer
            </button>
        </div>
        <div>{flash()}</div>
    </div>)
}

function KnownPeerList() {
    const [profiles, profileOpt] = createResource(fetchCount, ViewAllProfiles)
    onMount(refetch)
    return (<div class="m-8 overflow-scroll no-scrollbar">
        <Show when={profiles()?.length > 0}
            fallback={<div class="h-[100%] w-[100%] flex place-content-center place-items-center">🦗🍃 It's lonely here...</div>}
        >
            <For each={profiles()}>{(profile, i) =>
                <div class="flex gap-4 place-center">
                    <div><Image pubkey={profile.pubkey}/></div>
                    <div class="flex flex-col place-content-center pb-7">                    
                        <div>Alias: {profile.alias}</div>
                        <div>Public key: {profile.pubkey}</div>
                    </div>
                </div>
            }</For>
        </Show>
    </div>)
}

type Profile = { pubkey: string }

function Image({ pubkey }: Profile) {

    const svgURI = createMemo(() => "data:image/svg+xml;utf8," + encodeURIComponent(minidenticon(pubkey)))

    return (
        <div class="w-14 h-14 border rounded-full flex items-center justify-center">
            <img height={80} width={80} src={svgURI()} alt="profile picture" />
            <div ></div>
        </div>
    )
}