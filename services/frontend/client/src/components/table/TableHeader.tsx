import type { Clickable, HasChildren, HasClassname } from "../common";
import TableRow from "./TableRow";

export default function TableHeader(props: HasChildren & Clickable & HasClassname){
    return <TableRow {...props} className=" order-1"></TableRow>
}