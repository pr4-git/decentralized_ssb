import { MainView } from "./MainView";
import {Sync} from "./Sync";
import {Following} from "./Following";
import { Show, Switch, Match, createSignal, createResource, onMount } from 'solid-js';
import { IsServer } from "../wailsjs/go/main/App"
import serverIcon from './assets/server.svg';

export function App() {
    const [activeButtonIndex, setActiveButtonIndex] = createSignal(0);
    const defaultButtonStyle = "bg-transparent hover:bg-blue-500 text-white font-semibold hover:text-white py-2 px-4 border border-blue-500 hover:border-transparent rounded"
    const selectedButtonStyle = "bg-blue-600 text-white font-semibold hover:text-white py-2 px-4 border border-blue-500 hover:border-transparent rounded"

    

    const [isServer, { mutate, refetch }] = createResource(IsServer)
    onMount(refetch)

    return (
        <div class="w-screen h-screen grid grid-cols-[25%_75%]">

            <div class="col-start-1 col-span-1 border border-l-0 border-t-0 border-b-0 sticky left-0 flex flex-col justify-between">
                <nav class="flex flex-col gap-5 px-3 py-5 text-lg">
                    <button class={activeButtonIndex() === 0 ? selectedButtonStyle : defaultButtonStyle}
                        onclick={() => setActiveButtonIndex(0)}
                    >
                        Home
                    </button>
                    <button class={activeButtonIndex() === 1 ? selectedButtonStyle : defaultButtonStyle}
                        onclick={() => setActiveButtonIndex(1)}
                    >
                        Sync
                    </button>
                    <button class={activeButtonIndex() === 2 ? selectedButtonStyle : defaultButtonStyle}
                        onclick={() => setActiveButtonIndex(2)}
                    >
                        AddressBook
                    </button>
                </nav>

                <Show when={isServer()}>
                    <div class="py-2 flex flex-col my-20 justify-around">
                        <div class="flex flex-row place-content-center">
                            <div class="border rounded-full p-3">
                                <img src={serverIcon} color="white" width="40" height="40"></img>
                            </div>
                        </div>
                        <div class="m-2 text-center">You are the server for this instance.</div>
                    </div>
                </Show>
            </div>

            <div class="col-start-2 col-span-1 flex justify-center overflow-scroll no-scrollbar">
                <Switch fallback={<div> Something has gone terribly wrong!!</div>}>
                    <Match when={activeButtonIndex() === 0}>
                        <MainView />
                    </Match>
                    <Match when={activeButtonIndex() === 1}>
                        <Sync/>
                    </Match>
                    <Match when={activeButtonIndex() === 2}>
                        <Following/>
                    </Match>
                </Switch>
            </div>
        </div>
    )
}