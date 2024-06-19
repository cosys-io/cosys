import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@components/ui/card";
import { Label } from "@components/ui/label";
import { Input } from "@components/ui/input";
import { Button } from "@components/ui/button";
import { ClientActionFunctionArgs, Form, Link, redirect } from "@remix-run/react";

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
    const body = await request.formData();
    console.log(body.get("email"));
    return redirect("/");
};

export default function Login() {
    return (
        <div className="flex items-center justify-center h-screen w-screen">
            <Card className="mx-auto max-w-sm">
                <CardHeader>
                    <CardTitle className="text-2xl">Login</CardTitle>
                    <CardDescription>Enter your email below to login to your account</CardDescription>
                </CardHeader>
                <CardContent>
                    <Form method="POST">
                        <div className="grid gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="email">Email</Label>
                                <Input id="email" name="email" type="email" placeholder="m@example.com" required />
                            </div>
                            <div className="grid gap-2">
                                <div className="flex items-center">
                                    <Label htmlFor="password">Password</Label>
                                </div>
                                <Input id="password" name="password" type="password" required />
                            </div>
                            <Button type="submit" className="w-full">
                                Login
                            </Button>
                        </div>
                        <div className="mt-4 text-center text-sm">
                            Forgot your password?{" "}
                            <Link to={"/reset-password"} className="ml-auto inline-block text-sm underline">
                                Reset
                            </Link>
                        </div>
                    </Form>
                </CardContent>
            </Card>
        </div>
    );
}
