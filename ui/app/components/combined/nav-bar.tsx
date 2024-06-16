import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@components/ui/tooltip";
import { Link } from "@remix-run/react";
import { ChevronLeftIcon, ChevronRightIcon, Layout, SettingsIcon } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@components/ui/avatar";
import { Button } from "@components/ui/button";
import { useState } from "react";

const ICON_SIZE = "w-6 h-6";

export default function NavBar() {
    const [isNavCollapsed, setIsNavCollapsed] = useState(false);
    const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);

    return (
        <div
            className={`flex flex-col  px-4 py-8 bg-white border-r dark:bg-gray-800 dark:border-gray-600 transition-all duration-300 ${
                isNavCollapsed ? "w-28" : "w-52"
            }`}
        >
            <div className="flex items-center justify-start">
                <Link to="/" className={`transition-all duration-300 ${!isNavCollapsed ? "ml-4" : ""}`}>
                    <h2 className="text-2xl font-semibold text-gray-800 dark:text-white">Cosys</h2>
                </Link>
            </div>
            <div className="flex flex-col justify-between flex-1 mt-6">
                <nav>
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Link
                                    to="/resources"
                                    className={`flex items-center px-4 py-2 mt-5 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:text-gray-700 dark:hover:text-gray-200 ${
                                        isNavCollapsed ? "justify-center" : ""
                                    }`}
                                >
                                    <Layout className={ICON_SIZE} />
                                    {!isNavCollapsed && <span className="mx-4 font-medium">Resources</span>}
                                </Link>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>Resources</p>
                            </TooltipContent>
                        </Tooltip>
                        <hr className="my-6 dark:border-gray-600" />
                        <Tooltip>
                            <Link
                                to="#"
                                className={`flex items-center px-4 py-2 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:text-gray-700 dark:hover:text-gray-200 ${
                                    isNavCollapsed ? "justify-center" : ""
                                }`}
                            >
                                <SettingsIcon className={ICON_SIZE} />
                                {!isNavCollapsed && <span className="mx-4 font-medium">Settings</span>}
                            </Link>
                        </Tooltip>
                    </TooltipProvider>
                </nav>
                <div className="flex items-center">
                    <button
                        className="flex items-center px-4 py-2 mt-5 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700 hover:text-gray-700 dark:hover:text-gray-200 cursor-pointer"
                        onClick={() => setIsProfileMenuOpen(!isProfileMenuOpen)}
                    >
                        <Avatar className="w-8 h-8">
                            <AvatarImage src="https://github.com/shadcn.png" />
                            <AvatarFallback>JP</AvatarFallback>
                        </Avatar>
                        {!isNavCollapsed && <span className="mx-4 font-medium">Profile</span>}

                        {isProfileMenuOpen && (
                            <div className="absolute right-4 mt-2 w-48 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-gray-800">
                                <div className="py-1">
                                    <Link
                                        to="/profile"
                                        className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                                    >
                                        Profile
                                    </Link>
                                    <Link
                                        to="#"
                                        className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                                    >
                                        Logout
                                    </Link>
                                </div>
                            </div>
                        )}
                    </button>
                    <Button
                        className="translate-y-1/4 bg-white transition-all duration-300"
                        variant="ghost"
                        size="icon"
                        onClick={() => setIsNavCollapsed(!isNavCollapsed)}
                    >
                        {isNavCollapsed ? (
                            <ChevronRightIcon className={ICON_SIZE} />
                        ) : (
                            <ChevronLeftIcon className={ICON_SIZE} />
                        )}
                    </Button>
                </div>
            </div>
        </div>
    );
}
