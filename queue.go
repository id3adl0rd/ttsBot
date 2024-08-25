package main

/*func (b *Bot) GetSound() string {
	b.mutex.RLock()
	if len(b.queue) > 0 {
		return b.queue[0]
	}
	b.mutex.RUnlock()

	return ""
}

func (b *Bot) AddToQueue(message string) {
	b.mutex.Lock()
	b.queue = append(b.queue, message)
	b.mutex.Unlock()
}

func (b *Bot) QueueRemoveFisrt() {
	fmt.Println("locking")
	b.mutex.Lock()
	if len(b.queue) != 0 {
		b.queue = b.queue[1:]
	}
	fmt.Println("deleting")
	b.mutex.Unlock()
}

func (b *Bot) QueueRemoveIndex(k int) {
	b.mutex.Lock()
	if len(b.queue) != 0 && k <= len(b.queue) {
		b.queue = append(b.queue[:k], b.queue[k+1:]...)
	}
	b.mutex.Unlock()
}

func (b *Bot) QueueRemoveLast() {
	b.mutex.Lock()
	if len(b.queue) != 0 {
		b.queue = append(b.queue[:len(b.queue)-1], b.queue[len(b.queue):]...)
	}
	b.mutex.Unlock()
}

func (b *Bot) QueueClean() {
	b.mutex.Lock()
	b.queue = b.queue[:1]
	b.mutex.Unlock()
}

func (b *Bot) QueueRemove() {
	b.mutex.Lock()
	b.queue = []string{}
	b.mutex.Unlock()
}*/
