import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@components/ui/card";
import { Label } from "@components/ui/label";
import { Input } from "@components/ui/input";
import { Button } from "@components/ui/button";
import { Link } from "@remix-run/react";

export default function ResetPassword() {
    return (
        <div className="h-svh flex items-center justify-center">
            <Card className="mx-auto max-w-sm">
                <CardHeader>
                    <CardTitle className="text-2xl">Reset Password</CardTitle>
                    <CardDescription>Enter your email to reset your password</CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="grid gap-4">
                        <div className="grid gap-2">
                            <Label htmlFor="email">Email</Label>
                            <Input id="email" type="email" placeholder="m@example.com" required />
                        </div>
                        <Button type="submit" className="w-full">
                            Send Verification Code
                        </Button>
                    </div>
                    <div className="mt-4 text-center text-sm">
                        Remember your password?{" "}
                        <Link to="/login" className="underline">
                            Login
                        </Link>
                    </div>
                    <div className="mt-4 text-center text-sm text-gray-500 dark:text-gray-400">
                        An email has been sent to your email address with instructions to reset your password.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
