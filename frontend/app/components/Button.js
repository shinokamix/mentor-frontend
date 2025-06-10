'use client'

import React from 'react'
import { useRouter } from 'next/navigation';

const Button = ({text, href, alt=false, className=''}) => {
    const router = useRouter();

    const handleClick = () => {
        router.push(href);
    };

    const altStyle = alt ? 'bg-white text-black' : 'bg-blue-400 text-white';

    return (
        <div className={`cursor-pointer py-2 px-5 rounded-2xl ${altStyle}`} onClick={handleClick}>
            <p>{text}</p>
        </div>
    )
}

export default Button