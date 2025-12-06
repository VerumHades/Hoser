
export interface Rental {
    id: string;
    title: string;
    description: string;
}

interface RentalDisplayProps {
    rental: Rental,
    onShouldClose?: () => void
}

export default function UserRentalDisplay({ rental }: RentalDisplayProps) {
    return (
        <div className="w-full h-full p-6  shadow-md rounded-lg">
            <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                <label>
                    {rental.title}
                </label>
            </h1>
            <p className="mb-2 text-gray-700 dark:text-gray-300">
                <label>
                    {rental.description}
                </label>
            </p>
        </div>
    );
};
