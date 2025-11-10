import {  LoaderCircle} from "lucide-react";

export default function LoadingIcon(){
    return <div className="w-full h-full flex flex-row justify-center items-center">
        <LoaderCircle className="text-center w-12 h-12 text-blue-500 animate-spin"/>
    </div>
}