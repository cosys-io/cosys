import NavBar from "@components/combined/nav-bar";
import { Input } from "@components/ui/input";
import { Button } from "@components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@components/ui/table";
import { Checkbox } from "@components/ui/checkbox";
import { PencilIcon, TrashIcon } from "lucide-react";

export default function Resources() {
    return (
        <div className="flex h-screen w-screen bg-gray-100">
            <NavBar />
            <div className="flex flex-col w-48 px-4 py-8 bg-white border-r dark:bg-gray-800 dark:border-gray-600">
                <h2 className="text-2xl font-semibold text-gray-800 dark:text-white">Resource</h2>
                <div className="flex flex-col justify-between flex-1 mt-6">
                    <nav>
                        <a
                            href="/"
                            className="w-full px-4 py-2 mt-2 text-sm font-medium text-gray-600 dark:text-gray-400 dark:hover:text-white hover:bg-gray-200 dark:hover:bg-gray-700"
                        >
                            Access Link
                        </a>
                    </nav>
                </div>
            </div>
            <main className="flex-1 overflow-x-hidden overflow-y-auto bg-gray-200">
                <div className="container mx-auto px-6 py-8">
                    <h3 className="text-gray-700 text-3xl font-medium">Content Manager</h3>
                    <div className="mt-4 flex flex-wrap">
                        <div className="w-full mb-4 lg:mb-0">
                            <div className="overflow-hidden bg-white shadow-lg sm:rounded-lg">
                                <div className="px-6 py-4">
                                    <div className="flex items-center justify-between">
                                        <h3 className="text-xl font-semibold text-gray-800 dark:text-white">
                                            Access Link
                                        </h3>
                                        <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
                                            5,334 entries found
                                        </span>
                                    </div>
                                    <div className="mt-4">
                                        <div className="flex justify-between">
                                            <Input placeholder="Search" className="w-full" />
                                            <Button className="ml-4">Filters</Button>
                                        </div>
                                        <div className="mt-6 w-full">
                                            <Table>
                                                <TableHeader>
                                                    <TableRow>
                                                        <TableHead className="w-20">
                                                            <Checkbox id="check-all" />
                                                        </TableHead>
                                                        <TableHead>ID</TableHead>
                                                        <TableHead>CUSTOMER</TableHead>
                                                        <TableHead>SALESPERSON</TableHead>
                                                        <TableHead>ACTIVE</TableHead>
                                                        <TableHead>SLUG</TableHead>
                                                        <TableHead>ACTIONS</TableHead>
                                                    </TableRow>
                                                </TableHeader>
                                                <TableBody>
                                                    <TableRow>
                                                        <TableCell>
                                                            <Checkbox id="entry-5407" />
                                                        </TableCell>
                                                        <TableCell className="font-medium">5407</TableCell>
                                                        <TableCell>mia collection</TableCell>
                                                        <TableCell>Luna Pinkdose</TableCell>
                                                        <TableCell>true</TableCell>
                                                        <TableCell>36ff2368-93ad-43eb-adfb-f66876014b80</TableCell>
                                                        <TableCell>
                                                            <div className="flex justify-center gap-2">
                                                                <PencilIcon className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                                                                <TrashIcon className="w-5 h-5 text-gray-600 dark:text-gray-400" />
                                                            </div>
                                                        </TableCell>
                                                    </TableRow>
                                                </TableBody>
                                            </Table>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}
