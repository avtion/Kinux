import { createAvatar } from '@dicebear/avatars'
import * as AvatarsSprites from '@dicebear/avatars-avataaars-sprites'
import * as IntSprites from '@dicebear/avatars-initials-sprites'

// 头像生成
export const ProfileAvatarCreator = (seed: string) => {
    return createAvatar(AvatarsSprites, {
        seed: seed,
        dataUri: true,
        skin: ['light'],
    })
}

export const IntCreator = (seed: string, color: string = '#1D4ED8') => {
    return createAvatar(IntSprites, {
        seed: seed,
        dataUri: true,
        background: color,
    })
}