import { Input } from "./ui/input";
import { Label } from "./ui/label";

interface FormInputProps extends React.ComponentProps<"input"> {
  label?: string;
  error?: string | string[] | undefined;
}

export const FormInput = ({ label, error, ...props }: FormInputProps) => {
  return (
    <div className="space-y-2">
      {label && <Label>{label}</Label>}
      <Input {...props} className="bg-transparent" />
      {error && <p className="text-destructive text-sm">{error}</p>}
    </div>
  );
};
