import React from 'react'
import Button from './Button'

const Header = () => {
  return (
    <div className=''>
        <div className='flex my-3'>
            <div className='ml-3'>
                <Button text='Home' href={'/'} />
            </div>
            <div className='flex gap-3 ml-auto mr-3'>
                <Button text='Sign In' href={'/login'} alt={true} />
                <Button text='Sign Up' href={'/register'} alt={false} />
            </div>
        </div>
    </div>
  )
}

export default Header