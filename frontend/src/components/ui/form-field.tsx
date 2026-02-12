import type { ComponentProps } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface FormFieldProps extends ComponentProps<typeof Input> {
  label?: string;
}

export function FormField({ label, id, ...props }: FormFieldProps) {
  return (
    <div className="space-y-1">
      {label && <Label htmlFor={id}>{label}</Label>}
      <Input id={id} {...props} />
    </div>
  );
}
