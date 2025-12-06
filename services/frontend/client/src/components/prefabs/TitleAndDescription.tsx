
interface TitleAndDescriptionProps {
    title: string,
    titleClassname?: string,
    description: string
    descriptionClassname?: string
    className?: string
}

export function TitleAndDescription({ title, description, titleClassname, descriptionClassname, className }: TitleAndDescriptionProps) {
    return  <div className={"flex flex-col " + className}>
            <label className={"text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100 " + titleClassname}>
                {title}
            </label>
            <label className={"mb-2 text-gray-700 dark:text-gray-300 " + descriptionClassname }>
                {description}
            </label>
        </div>
}