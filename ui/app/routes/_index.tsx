import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
    return [
        { title: "Cosys Content Management UI" },
        { name: "description", content: "Hope you have fun building your backend!" },
    ];
};
import NavBar from "@components/combined/nav-bar";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@components/ui/card";

export default function Index() {
    return (
        <div className="flex h-screen w-screen bg-gray-100">
            <NavBar />
            <div className="w-full flex justify-center items-center">
                <Card className="max-w-2xl">
                    <CardHeader>
                        <CardTitle>Welcome to Cosys</CardTitle>
                        <CardDescription>We are a microservice framework</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <p>We hope you enjoy using the framework</p>
                    </CardContent>
                    <CardFooter>
                        <p className="text-right w-full text-gray-400">By Adwin and Yu Han</p>
                    </CardFooter>
                </Card>
            </div>
        </div>
    );
}
