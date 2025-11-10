// src/components/DashboardApp.tsx
import Search from "../components/Search";

interface ElementListProps<T> {
    element: T | undefined,
    children?: React.ReactNode,

    itemBuilder: (item: T) => React.ReactElement,
    itemViewBuilder: (item: T) => React.ReactElement,
    
    endpoint: string,
    queryBuilder: (query: string) => {},
}

export default function ElementList<T>({queryBuilder, itemBuilder, itemViewBuilder, element, endpoint}: ElementListProps<T>) {
    if (element) return itemViewBuilder(element)

    const builder = (element: T) => {
        return itemBuilder(element)
    }

    return <div className="w-full h-full flex flex-col items-center">
        <Search itemBuilder={builder} queryBuilder={queryBuilder} endpoint={endpoint}>

        </Search>
    </div>
}
