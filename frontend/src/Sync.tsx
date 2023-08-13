import { For, Show, createResource, createSignal, onMount } from "solid-js"
import { GetMyPublickey, IsServer, JoinGossipNet } from "../wailsjs/go/main/App"
import connection from "./assets/connection.svg";
import connectionNone from "./assets/connectionNone.svg";
import syncProgress from "./assets/syncProgress.svg";
import { Icon } from "./Icon";

export function Sync() {
    const [isServer, { mutate, refetch }] = createResource(IsServer)
    onMount(() => {
        refetch()
    })


    return (<div>
        <Show when={isServer()}
            fallback={<ClientSync />}
        >
            <ServerSync />
        </Show>

    </div>)
}

function ServerSync() {
    const [pk, { mutate, refetch }] = createResource(GetMyPublickey)
    const [copyNotif, setcopyNotif] = createSignal("📋Copy to clipboard")


    const notifyCopy = () => {
        let msg = "👌Copied to clipboard"
        let before = copyNotif()
        setcopyNotif(msg)
        setTimeout(() => setcopyNotif(before), 2000)
    }

    onMount(refetch)
    const CopyPkToClipboard = () => {
        navigator.clipboard.writeText(pk() || "ssb-ng: null pubkey")
        notifyCopy()
    }

    return (<div class="flex flex-col place-content-center h-screen">
        <Icon iconSrc={connection} msg="Serving all peers on 127.0.0.1:8008" />
        <div class="text-center">
            <div>
                Please copy your public key for authenticating other peers.
            </div>
            <div class="text-white text-xs py-2 mt-3 bg-gray-800">
                {pk()}
            </div>
        </div>
        <div class="flex place-content-center mt-2 font-bold">
            <button class="border rounded-full p-1 px-3 bg-blue-500 border-blue-500 hover:bg-blue-400"
                onClick={CopyPkToClipboard}>
                {copyNotif()}
            </button>
        </div>
    </div>)
}

let syncingPeers = []

function ClientSync() {
    const [pk, { mutate, refetch }] = createResource(GetMyPublickey)
    const [copyNotif, setcopyNotif] = createSignal("📋Copy to clipboard")


    const notifyCopy = () => {
        let msg = "👌Copied to clipboard"
        let before = copyNotif()
        setcopyNotif(msg)
        setTimeout(() => setcopyNotif(before), 2000)
    }

    onMount(refetch)
    const CopyPkToClipboard = () => {
        navigator.clipboard.writeText(pk() || "ssb-ng: null pubkey")
        notifyCopy()
    }


    const buttonEnabled = "border rounded-full p-1 px-3 bg-blue-500 border-blue-500 hover:bg-blue-400"
    const buttonDisabled = "border rounded-full p-1 px-3 text-gray-600 bg-transparent border-gray-600"

    const [pkfield, setPkfield] = createSignal("...")
    const [syncing, setSyncing] = createSignal(false)

    const StartSync = () => {
        JoinGossipNet(pkfield())
            .then(() => setSyncing(true))
            .catch(err => console.log(err))
        syncingPeers.push(pkfield())
    }

    onMount(() => { (syncingPeers.length == 0) ? setSyncing(false) : setSyncing(true) })

    return (<div class="flex flex-col place-content-center h-screen">
        <Show when={syncing()}
            fallback={
                <Icon iconSrc={connectionNone}
                    msg={"Not connected to any gossip network"} />
            }
        >
            <Icon iconSrc={syncProgress}
                msg={"Connected and Syncing"} />
        </Show>

        <div class="text-center">
            <div>
                Please Enter the public key of server to connect.
            </div>
            <input class="text-white text-center w-full py-2 mt-3 h-8 bg-gray-800"
                placeholder="Public key of server"
                onInput={(e) => { setPkfield(e.target.value) }}
            >
            </input>
        </div>
        <div class="flex place-content-center mt-2 font-bold">
            <button class={syncing() ? buttonDisabled : buttonEnabled}
                onClick={StartSync}>
                Sync
            </button>
        </div>
        <div class="flex flex-col place-content-center mt-5 ">
            <div class="text-sm text-center">Your Public key is:</div>
            <div class="text-xs mt-2 bg-gray-800 h-8 flex place-items-center place-content-center">
            {pk()}
            </div>
            <button class="mx-12 mt-3 border rounded-full p-1 px-3 bg-blue-500 border-blue-500 hover:bg-blue-400"
                onClick={CopyPkToClipboard}>
                {copyNotif()}
            </button>
        </div>
    </div>)
}


