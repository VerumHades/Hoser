// src/components/AccountPage.tsx
import EditableText from "../components/EditableText";

export interface Rental {
    ID: string;
    Title: string;
    Description: string;
}

interface RentalDisplayProps {
    rental: Rental,
    onShouldClose?: () => void
}

export default function UserRentalDisplay({ rental }: RentalDisplayProps) {
    return (
        <div className="w-full h-full p-6  shadow-md rounded-lg">
            <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                <EditableText
                    text={rental.Title}
                    onChange={() => {}}>

                </EditableText>
            </h1>
            <p className="mb-2 text-gray-700 dark:text-gray-300">
                <EditableText
                    text={rental.Description}
                    onChange={() => {}}>

                </EditableText>
            </p>
        </div>
    );
};
